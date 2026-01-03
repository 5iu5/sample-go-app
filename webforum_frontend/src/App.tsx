import { useState } from 'react'

import './App.css'

import ForumHome from "./pages/ForumHome"
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';

function App() {


  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<ForumHome />} />
        <Route path="/forum" element={<ForumHome />} />
      </Routes>
    </BrowserRouter>
    
  )
}

export default App
