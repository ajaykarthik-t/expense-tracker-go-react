import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Register from './pages/Register';
import Login from './pages/Login';
import Transactions from './pages/Transaction';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import BackgroundLayout from './components/BackgroundLayout';
import './App.css';

function App() {
  return (
    <Router>
      <BackgroundLayout>
        <Navbar />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route path="/transactions" element={<Transactions />} />
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </BackgroundLayout>
    </Router>
  );
}

export default App;
