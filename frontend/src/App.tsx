import React from 'react'
import { Routes, Route, Navigate } from "react-router-dom";
import MainPage from "./pages/MainPage";
import Register from './pages/Register';

const App: React.FC = () => (
  <Routes>
    <Route path='*' element={<Navigate to='/' replace />} />
    <Route path='/' element={<MainPage />} />
    <Route path='/register' element={<Register />} />
  </Routes>
);


export default App
