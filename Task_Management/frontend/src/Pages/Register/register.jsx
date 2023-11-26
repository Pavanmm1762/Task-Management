import React, { useState } from "react";
import Card from "../../components/Card/card";
import "./register.css";
import GoogleIcon from "@mui/icons-material/Google";
import FacebookIcon from "@mui/icons-material/Facebook";
import TwitterIcon from "@mui/icons-material/Twitter";
import axios from "axios";
import { Link } from 'react-router-dom';

const RegisterForm = ({ setIsLoggedIn }) => {

  const [signupdata, setSignUpData] = useState({ username: "", email: "", password: "" });
  const [errorMessages, setErrorMessages] = useState({});

  const errors = {
    username: "Invalid username",
    password: "Invalid password",
    email: "Invalid email",
    noEmail: "Please enter your Email",
    noUsername: "Please enter your username",
    noPassword: "Please enter your password",

  };

  const handleLogin = async (e) => {
    e.preventDefault();


    if (!signupdata.username) {
      // Username input is empty
      setErrorMessages({ name: "noUsername", message: errors.noUsername });
      return;
    }

    else if (!signupdata.email) {
      // Email input is empty
      setErrorMessages({ name: "noEmail", message: errors.noEmail });
      return;
    }

    else if (!signupdata.password) {
      // Password input is empty
      setErrorMessages({ name: "noPassword", message: errors.noPassword });
      return;
    }

    else {
      try {
        const response = await axios.post("http://localhost:8080/auth/register", signupdata);
        console.log(response.data);
      } catch (error) {
        console.error("Error during register:", error.response.data);
      }
    }
  };

  // Render error messages
  const renderErrorMsg = (name) =>
    name === errorMessages.name && (
      <p className="error_msg">{errorMessages.message}</p>
    );

  return (
    <div className="block1">
      <Card>
        <h1 className="title">Sign Up</h1>
        <p className="subtitle">
          Please register using your email and create a password!
        </p>
        <form onSubmit={handleLogin}>
          <div className="inputs_container">
            <input
              type="text"
              placeholder="Username"
              onChange={(e) => setSignUpData({ ...signupdata, username: e.target.value })}
            />
            {renderErrorMsg("username")}
            {renderErrorMsg("noUsername")}
            <input
              type="email"
              placeholder="Email"
              onChange={(e) => setSignUpData({ ...signupdata, email: e.target.value })}
            />
            {renderErrorMsg("email")}
            {renderErrorMsg("noEmail")}
            <input
              type="password"
              placeholder="Password"
              onChange={(e) => setSignUpData({ ...signupdata, password: e.target.value })}
            />
            {renderErrorMsg("password")}
            {renderErrorMsg("noPassword")}
          </div>
          <input type="submit" value="Register" className="login_button" />
        </form>
        <div className="link_container">
          <span className="small">Already have an account? <Link to="/">Login</Link> </span>
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

export default RegisterForm;