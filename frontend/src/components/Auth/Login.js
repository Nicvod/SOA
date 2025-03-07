import React, { useState } from 'react';
import authService from '../../services/auth';

const Login = (props) => {
  const [credentials, setCredentials] = useState({ login: '', password: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await authService.login(credentials);
      props.setIsAuthenticated(true);
      window.location.href = '/#/profile';
    } catch (error) {
      console.error('Login failed', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Login"
        value={credentials.login}
        onChange={(e) => setCredentials({ ...credentials, login: e.target.value })}
      />
      <input
        type="password"
        placeholder="Password"
        value={credentials.password}
        onChange={(e) => setCredentials({ ...credentials, password: e.target.value })}
      />
      <button type="submit">Login</button>
    </form>
  );
};

export default Login;