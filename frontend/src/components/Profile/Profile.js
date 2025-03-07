import React, { useEffect, useState } from 'react';
import profileService from '../../services/profile';

const Profile = () => {
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const data = await profileService.getProfile();
        setProfile(data);
      } catch (error) {
        console.error('Failed to fetch profile', error);
      }
    };
    fetchProfile();
  }, []);

  if (!profile) return <div>Loading...</div>;

  return (
    <div>
      <h1>Profile</h1>
      <p>Login: {profile.login}</p>
      <p>Email: {profile.email}</p>
      <p>First Name: {profile.first_name}</p>
      <p>Last Name: {profile.last_name}</p>
      <p>Birth Date: {profile.birth_date}</p>
      <p>Phone Number: {profile.phone_number}</p>
      <a href="/#/profile/edit">Edit Profile</a>
    </div>
  );
};

export default Profile;