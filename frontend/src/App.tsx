import React from 'react'
import { Routes, Route, Navigate } from "react-router-dom";
import MainPage from "./pages/MainPage";


const App: React.FC = () => (
  <Routes>
    <Route path='*' element={<Navigate to='/' replace />} />
    <Route path='/' element={<MainPage />} />
  </Routes>
);


export default App
