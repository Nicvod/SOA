from dotenv import load_dotenv
import os
import pytest
import requests
from typing import List, Dict, Any

load_dotenv('.test.env')


class TestPostAPI:
    BASE_URL: str = os.getenv("TEST_API_BASE_URL", "http://localhost")
    LOGIN: str = os.getenv("TEST_API_LOGIN")
    PASSWORD: str = os.getenv("TEST_API_PASSWORD")
    
    @pytest.fixture(scope="class")
    def auth_token(self) -> str:
        if not self.LOGIN or not self.PASSWORD:
            raise ValueError("Необходимо установить API_LOGIN и API_PASSWORD в переменных окружения")
        
        auth_url = f"{self.BASE_URL}/api/v1/authenticate"
        auth_data = {
            "login": self.LOGIN,
            "password": self.PASSWORD
        }
        
        response = requests.post(auth_url, json=auth_data)
        assert response.status_code == 200, "Ошибка аутентификации"
        
        token = response.json().get("access_token")
        assert token is not None, "Токен не получен"
        
        return token
    
    @pytest.fixture(scope="class")
    def created_posts(self) -> List[str]:
        return []
    
    def test_create_post(self, auth_token: str, created_posts: List[str]):
        url = f"{self.BASE_URL}/api/v1/posts"
        headers = {"Authorization": f"Bearer {auth_token}"}
        post_data = {
            "title": "Test Post",
            "description": "This is a test post",
            "is_private": False,
            "tags": ["test", "integration"]
        }
        
        response = requests.post(url, json=post_data, headers=headers)
        assert response.status_code == 201, "Ошибка создания поста"
        
        post = response.json()
        assert "id" in post, "ID поста не возвращен"
        assert post["title"] == post_data["title"], "Некорректный заголовок поста"
        assert post["description"] == post_data["description"], "Некорректное описание поста"
        
        created_posts.append(post["id"])
    
    def test_list_posts(self, auth_token: str):
        url = f"{self.BASE_URL}/api/v1/posts"
        headers = {"Authorization": f"Bearer {auth_token}"}
        params = {"page": 0, "page_size": 10}
        
        response = requests.get(url, headers=headers, params=params)
        assert response.status_code == 200, "Ошибка получения списка постов"
        
        data = response.json()
        assert "posts" in data, "Список постов не возвращен"
        assert isinstance(data["posts"], list), "Посты должны быть списком"
        assert "total_count" in data, "Общее количество постов не возвращено"
    
    def test_get_post(self, auth_token: str, created_posts: List[str]):
        if not created_posts:
            pytest.skip("Нет созданных постов для тестирования")
        
        post_id = created_posts[0]
        url = f"{self.BASE_URL}/api/v1/posts/{post_id}"
        headers = {"Authorization": f"Bearer {auth_token}"}
        
        response = requests.get(url, headers=headers)
        assert response.status_code == 200, "Ошибка получения поста"
        
        post = response.json()
        assert post["id"] == post_id, "Некорректный ID поста"
        assert "title" in post, "Заголовок поста не возвращен"
        assert "description" in post, "Описание поста не возвращено"
    
    def test_update_post(self, auth_token: str, created_posts: List[str]):
        if not created_posts:
            pytest.skip("Нет созданных постов для тестирования")
        
        post_id = created_posts[0]
        url = f"{self.BASE_URL}/api/v1/posts/{post_id}"
        headers = {"Authorization": f"Bearer {auth_token}"}
        update_data = {
            "title": "Updated Test Post",
            "description": "This is an updated test post",
            "is_private": True,
            "tags": ["updated", "test"]
        }
        
        response = requests.put(url, json=update_data, headers=headers)
        assert response.status_code == 200, "Ошибка обновления поста"
        
        updated_post = response.json()
        assert updated_post["title"] == update_data["title"], "Заголовок не обновлен"
        assert updated_post["description"] == update_data["description"], "Описание не обновлено"
        assert updated_post["is_private"] == update_data["is_private"], "Статус приватности не обновлен"
        assert set(updated_post["tags"]) == set(update_data["tags"]), "Теги не обновлены"
    
    def test_delete_post(self, auth_token: str, created_posts: List[str]):
        if not created_posts:
            pytest.skip("Нет созданных постов для тестирования")
        
        for post_id in created_posts.copy():
            url = f"{self.BASE_URL}/api/v1/posts/{post_id}"
            headers = {"Authorization": f"Bearer {auth_token}"}
            
            response = requests.delete(url, headers=headers)
            assert response.status_code == 204, f"Ошибка удаления поста {post_id}"
            
            get_response = requests.get(url, headers=headers)
            assert get_response.status_code not in (200, 201), f"Пост {post_id} не был удален"
            
            created_posts.remove(post_id)
        
        assert len(created_posts) == 0, "Не все посты были удалены"