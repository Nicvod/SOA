package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Nicvod/SOA/postService/internal/models"

	post_proto "github.com/Nicvod/SOA/postService/post_proto"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(ctx context.Context, post *post_proto.CreatePostRequest, creatorID string) (*post_proto.PostResponse, error) {
	var resp post_proto.PostResponse
	now := time.Now()
	postID := uuid.New().String()

	query := `
		INSERT INTO posts (id, title, description, creator_id, created_at, updated_at, is_private, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, creator_id, created_at, updated_at, is_private, tags
	`

	err := r.db.QueryRowContext(ctx, query,
		postID,
		post.Title,
		post.Description,
		creatorID,
		now,
		now,
		post.IsPrivate,
		pq.Array(post.Tags),
	).Scan(
		&resp.Id,
		&resp.Title,
		&resp.Description,
		&resp.CreatorId,
		&now,
		&now,
		&resp.IsPrivate,
		pq.Array(&resp.Tags),
	)

	if err != nil {
		return nil, err
	}

	resp.CreatedAt = timestamppb.New(now)
	resp.UpdatedAt = timestamppb.New(now)
	return &resp, nil
}

func (r *PostRepository) GetPost(ctx context.Context, postID, userID string) (*post_proto.PostResponse, error) {
	var resp post_proto.PostResponse
	var createdAt, updatedAt time.Time

	query := `
		SELECT id, title, description, creator_id, created_at, updated_at, is_private, tags
		FROM posts
		WHERE id = $1 AND (is_private = FALSE OR creator_id=$2)
	`

	err := r.db.QueryRowContext(ctx, query, postID, userID).Scan(
		&resp.Id,
		&resp.Title,
		&resp.Description,
		&resp.CreatorId,
		&createdAt,
		&updatedAt,
		&resp.IsPrivate,
		pq.Array(&resp.Tags),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrPostNotFound
		}
		return nil, err
	}

	resp.CreatedAt = timestamppb.New(createdAt)
	resp.UpdatedAt = timestamppb.New(updatedAt)
	return &resp, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *post_proto.UpdatePostRequest, creatorID string) (*post_proto.PostResponse, error) {
	var resp post_proto.PostResponse
	var createdAt, updatedAt time.Time

	query := `
		UPDATE posts
		SET title = $2, description = $3, is_private = $4, tags = $5, updated_at = NOW()
		WHERE id = $1 AND creator_id = $6
		RETURNING id, title, description, creator_id, created_at, updated_at, is_private, tags
	`

	err := r.db.QueryRowContext(ctx, query,
		post.PostId,
		post.Title,
		post.Description,
		post.IsPrivate,
		pq.Array(post.Tags),
		creatorID,
	).Scan(
		&resp.Id,
		&resp.Title,
		&resp.Description,
		&resp.CreatorId,
		&createdAt,
		&updatedAt,
		&resp.IsPrivate,
		pq.Array(&resp.Tags),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrPostNotFound
		}
		return nil, err
	}

	resp.CreatedAt = timestamppb.New(createdAt)
	resp.UpdatedAt = timestamppb.New(updatedAt)
	return &resp, nil
}

func (r *PostRepository) DeletePost(ctx context.Context, postID, creatorID string) error {
	query := `
		DELETE FROM posts
		WHERE id = $1 AND creator_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, postID, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrPostNotFound
	}

	return nil
}

func (r *PostRepository) ListPosts(ctx context.Context, userID string, page, pageSize int32) (*post_proto.ListPostsResponse, error) {
	var response post_proto.ListPostsResponse
	var totalCount int32

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM posts WHERE is_private = FALSE OR creator_id = $1", userID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	offset := (page) * pageSize
	query := `
		SELECT id, title, description, creator_id, created_at, updated_at, is_private, tags
		FROM posts
		WHERE is_private = FALSE
		OR creator_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post post_proto.PostResponse
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.CreatorId,
			&createdAt,
			&updatedAt,
			&post.IsPrivate,
			pq.Array(&post.Tags),
		)
		if err != nil {
			return nil, err
		}

		post.CreatedAt = timestamppb.New(createdAt)
		post.UpdatedAt = timestamppb.New(updatedAt)
		response.Posts = append(response.Posts, &post)
	}

	response.TotalCount = totalCount
	response.Page = page
	response.PageSize = pageSize

	return &response, nil
}
