import React, { useState, useEffect } from 'react';
import profileService from '../../services/profile';

const EditProfile = () => {
  const [profile, setProfile] = useState({
    email: '',
    first_name: '',
    last_name: '',
    birth_date: '',
    phone_number: '',
  });

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

  const handleSubmit = async (e) => {
    e.preventDefault();

    const formattedProfile = {
      ...profile,
      birth_date: new Date(profile.birth_date).toISOString(),
    };
    try {
      await profileService.updateProfile(formattedProfile);
      window.location.href = '/#/profile';
    } catch (error) {
      console.error('Failed to update profile', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        placeholder="Email"
        value={profile.email}
        onChange={(e) => setProfile({ ...profile, email: e.target.value })}
      />
      <input
        type="text"
        placeholder="First Name"
        value={profile.first_name}
        onChange={(e) => setProfile({ ...profile, first_name: e.target.value })}
      />
      <input
        type="text"
        placeholder="Last Name"
        value={profile.last_name}
        onChange={(e) => setProfile({ ...profile, last_name: e.target.value })}
      />
      <input
        type="date"
        placeholder="Birth Date"
        value={profile.birth_date}
        onChange={(e) => setProfile({ ...profile, birth_date: e.target.value })}
      />
      <input
        type="text"
        placeholder="Phone Number"
        value={profile.phone_number}
        onChange={(e) => setProfile({ ...profile, phone_number: e.target.value })}
      />
      <button type="submit">Update Profile</button>
    </form>
  );
};

export default EditProfile;