import React, { useState } from 'react';
import authService from '../../services/auth';

const Register = (props) => {
  const [userData, setUserData] = useState({
    login: '',
    password: '',
    email: '',
    first_name: '',
    last_name: '',
    birth_date: '',
    phone_number: '',
  });

  const handleSubmit = async (e) => {
    e.preventDefault();

    const formattedUserData = {
      ...userData,
      birth_date: new Date(userData.birth_date).toISOString(),
    };
    try {
      await authService.register(formattedUserData);
      props.setIsAuthenticated(true);
      window.location.href = '/#/profile';
    } catch (error) {
      console.error('Registration failed', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Login"
        value={userData.login}
        onChange={(e) => setUserData({ ...userData, login: e.target.value })}
      />
      <input
        type="password"
        placeholder="Password"
        value={userData.password}
        onChange={(e) => setUserData({ ...userData, password: e.target.value })}
      />
      <input
        type="email"
        placeholder="Email"
        value={userData.email}
        onChange={(e) => setUserData({ ...userData, email: e.target.value })}
      />
      <input
        type="text"
        placeholder="First Name"
        value={userData.first_name}
        onChange={(e) => setUserData({ ...userData, first_name: e.target.value })}
      />
      <input
        type="text"
        placeholder="Last Name"
        value={userData.last_name}
        onChange={(e) => setUserData({ ...userData, last_name: e.target.value })}
      />
      <input
        type="date"
        placeholder="Birth Date"
        value={userData.birth_date}
        onChange={(e) => setUserData({ ...userData, birth_date: e.target.value })}
      />
      <input
        type="text"
        placeholder="Phone Number"
        value={userData.phone_number}
        onChange={(e) => setUserData({ ...userData, phone_number: e.target.value })}
      />
      <button type="submit">Register</button>
    </form>
  );
};

export default Register;