// App.js
import React, {useState}  from "react";
import './App.css';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';

import LoginForm from "./components/Login/login";
import RegisterForm from "./components/Register/register";
import Dashboard from "./components/Mainpage/pages/dashboard";
import Sidebar from "./components/Mainpage/sidebar";
import Project from "./components/Mainpage/pages/project";
import Tasks from "./components/Mainpage/pages/task";
import Report from "./components/Mainpage/pages/report";
import Users from "./components/Mainpage/pages/users";
import About from "./components/Mainpage/pages/about";
import Header from "./components/Mainpage/header";


const App = () => {
  const [isLoggedIn, setLoggedIn] = useState(false);

  return (
    <div>
      <Router>
        <Header/>
        <Routes>
          <Route
            path="/"
            element={isLoggedIn ? <Dashboard /> : <LoginForm onLogin={() => setLoggedIn(true)} />}
          />
          <Route path="/register" element={<RegisterForm />} />
          <Route
            path="/task_management/*"
            element={
            <React.Fragment>
              <Sidebar >
              <Routes>
                <Route path="/dashboard" element={<Dashboard />} />
                <Route path="/project" element={<Project />} />
                <Route path="/tasks" element={<Tasks />} />
                <Route path="/report" element={<Report />} />
                <Route path="/users" element={<Users />} />
                <Route path="/about" element={<About />} />
              </Routes>
              </Sidebar>
             </React.Fragment>
            }
          />
        </Routes>
      </Router>
    </div>
  );
};

 

export default App;
