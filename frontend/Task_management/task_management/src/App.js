// App.js
import React, { useState } from "react";
import axios from "axios";

function App() {
  const [signupData, setSignupData] = useState({ username: "", password: "" });
  const [loginData, setLoginData] = useState({ username: "", password: "" });

  const handleSignup = async () => {
    try {
      const response = await axios.post("http://localhost:8080/api/signup", signupData);
      console.log(response.data);
    } catch (error) {
      console.error("Error during signup:", error);
    }
  };

  const handleLogin = async () => {
    try {
      const response = await axios.post("http://localhost:8080/api/login", loginData);
      console.log(response.data);
    } catch (error) {
      console.error("Error during login:", error);
    }
  };

  return (
    <div>
      <h1>React, Golang Gin, and Cassandra Example</h1>

      <div>
        <h2>Signup</h2>
        <input
          type="text"
          placeholder="Username"
          onChange={(e) => setSignupData({ ...signupData, username: e.target.value })}
        />
        <input
          type="password"
          placeholder="Password"
          onChange={(e) => setSignupData({ ...signupData, password: e.target.value })}
        />
        <button onClick={handleSignup}>Signup</button>
      </div>

      <div>
        <h2>Login</h2>
        <input
          type="text"
          placeholder="Username"
          onChange={(e) => setLoginData({ ...loginData, username: e.target.value })}
        />
        <input
          type="password"
          placeholder="Password"
          onChange={(e) => setLoginData({ ...loginData, password: e.target.value })}
        />
        <button onClick={handleLogin}>Login</button>
      </div>
    </div>
  );
}

export default App;
