import React, { useState } from 'react';
import { HashRouter as Router, Route, Routes, Link } from 'react-router-dom';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import Profile from './components/Profile/Profile';
import EditProfile from './components/Profile/EditProfile';
import logout from './services/auth'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const handleLogout = () => {
    console.log(isAuthenticated, setIsAuthenticated)
    logout();
    setIsAuthenticated(false);
    console.log(isAuthenticated, setIsAuthenticated)
  };

  console.log(isAuthenticated, setIsAuthenticated)

  return (
    <Router>
      <nav>
        {!isAuthenticated && <Link to="/login">Login</Link>}
        {!isAuthenticated && <Link to="/register">Register</Link>}
        {isAuthenticated && <Link to="/profile">Profile</Link>}
        {isAuthenticated && <Link to="/profile/edit">Edit Profile</Link>}
        {!isAuthenticated && <button onClick={handleLogout}>Logout</button>}

      </nav>
      <Routes>
        <Route path="/login" element={<Login setIsAuthenticated={setIsAuthenticated} />} />
        <Route path="/register" element={<Register setIsAuthenticated={setIsAuthenticated}/>} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/profile/edit" element={<EditProfile />} />
      </Routes>
    </Router>
  );
}

export default App;