import axios from 'axios';

const API_URL = 'http://localhost/api/v1';

const register = async (userData) => {
  const response = await axios.post(`${API_URL}/register`, userData);
  if (response.data.access_token) {
    localStorage.setItem('access_token', response.data.access_token);
    localStorage.setItem('refresh_token', response.data.refresh_token);
  }
  return response.data;
};

const login = async (credentials) => {
  const response = await axios.post(`${API_URL}/authenticate`, credentials);
  if (response.data.access_token) {
    localStorage.setItem('access_token', response.data.access_token);
    localStorage.setItem('refresh_token', response.data.refresh_token);
  }
  return response.data;
};

const logout = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
};

const refreshToken = async () => {
  const refresh_token = localStorage.getItem('refresh_token');
  const response = await axios.post(`${API_URL}/refresh-token`, { refresh_token });
  if (response.data.access_token) {
    localStorage.setItem('access_token', response.data.access_token);
  }
  return response.data;
};

export default {
  register,
  login,
  logout,
  refreshToken,
};