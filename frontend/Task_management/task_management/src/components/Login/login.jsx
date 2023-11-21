import React, { useState } from "react";
import "./login.css";
import Card from "../Card/card";
import GoogleIcon from "@mui/icons-material/Google";
import FacebookIcon from "@mui/icons-material/Facebook";
import TwitterIcon from "@mui/icons-material/Twitter";
import axios from "axios";
import { Link } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';

const LoginForm = ({ onLogin }) => {

  const [loginData, setLoginData] = useState({ username: "", password: "" });
  const [errorMessages, setErrorMessages] = useState({});
  const navigate = useNavigate();

  const errors = {
    username: "Invalid username",
    password: "Invalid password",
    noUsername: "Please enter your username",
    noPassword: "Please enter your password",
  };

  const handleLogin = async (e) => {
    e.preventDefault();
  

    if (!loginData.username) {
      // Username input is empty
      setErrorMessages({ name: "noUsername", message: errors.noUsername });
      return;
    }

    else if (!loginData.password) {
      // Password input is empty
      setErrorMessages({ name: "noPassword", message: errors.noPassword });
      return;
    }

    else{
      try {
        const response = await axios.post("http://localhost:8080/api/login", loginData);
        onLogin();
        navigate('/task_management/dashboard');
        console.log(response.data);
      } catch (error) {
        console.error("Error during login:", error.response.data);
      }
    }
  };

  // Render error messages
  const renderErrorMsg = (name) =>
    name === errorMessages.name && (
      <p className="error_msg">{errorMessages.message}</p>
    );

  return (
    <div className="block">
      <Card>
        <h1 className="title">Sign In</h1>
        <p className="subtitle">
          Please log in using your username and password!
        </p>
        <form onSubmit={handleLogin}>
          <div className="inputs_container">
            <input
              type="text"
              placeholder="Username"
              onChange={(e) => setLoginData({ ...loginData, username: e.target.value })}
            />
            {renderErrorMsg("username")}
            {renderErrorMsg("noUsername")}
            <input
              type="password"
              placeholder="Password"
              onChange={(e) => setLoginData({ ...loginData, password: e.target.value })}
            />
            {renderErrorMsg("password")}
            {renderErrorMsg("noPassword")}
          </div>
          <input type="submit" value="Log In" className="login_button" />
        </form>
        <div className="link_container">
        <span className="small">New Here?  <Link to="/register">Register...</Link> </span>
        </div>
        <div className="icons">
          <GoogleIcon className="icon" />
          <FacebookIcon className="icon" />
          <TwitterIcon className="icon" />
        </div>
      </Card>
    </div>
  );
};

export default LoginForm;