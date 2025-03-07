import axios from 'axios';

const API_URL = 'http://localhost/api/v1';

const getProfile = async () => {
  const access_token = localStorage.getItem('access_token');
  const response = await axios.get(`${API_URL}/profile`, {
    headers: {
      Authorization: `Bearer ${access_token}`,
    },
  });
  return response.data;
};

const updateProfile = async (profileData) => {
  const access_token = localStorage.getItem('access_token');
  const response = await axios.put(`${API_URL}/profile`, profileData, {
    headers: {
      Authorization: `Bearer ${access_token}`,
    },
  });
  return response.data;
};

export default {
  getProfile,
  updateProfile,
};