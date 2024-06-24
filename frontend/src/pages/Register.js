import React, { useState } from 'react';

const Register = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState(''); // Added email state
  const [message, setMessage] = useState('');

  const register = () => {
    fetch('http://127.0.0.1:5000/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password, email }) // Include email in the request body
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        setMessage(data.error);
      } else {
        setMessage(data.message);
      }
    });
  };

  return (
    <div className="container">
      <h2>Register</h2>
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <input
        type="email" // Email input type
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button onClick={register}>Register</button>
      <div className="message">{message}</div>
    </div>
  );
};

export default Register;
