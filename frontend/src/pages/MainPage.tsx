import React, { useState } from "react";
import { useNavigate } from "react-router-dom"
import Layout from "../components/Layout";
import { backendClient } from "../data/backend";
import type { Dispatch, SetStateAction } from "react"

interface loginProp {
  stateHandler: Dispatch<SetStateAction<boolean>>
}

const Login: React.FC<loginProp> = ({ stateHandler }: loginProp) => {
  const [formData, setFormData] = useState({
    username: "",
    password: "",
  })

  const navigate = useNavigate();

  const handleLogin: React.FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault()

    try {
      const resp = await backendClient.post("/auth/login", formData)
      const { accessToken } = resp.data;
      localStorage.setItem("accessToken", accessToken);
      stateHandler(true)
    } catch (err) {
      if (err instanceof Error) {
        console.error(err.message)
      }

    }
  }

  const handleChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  return (
    <div className="login-form">
      <h1> Login to continue </h1>
      <form onSubmit={handleLogin}>
        <label>
          <p>Username</p>
          <input type="text" name="username" value={formData.username} onChange={handleChange} required />
        </label>
        <label>
          <p>Password</p>
          <input type="password" name="password" value={formData.password} onChange={handleChange} required />
        </label>
        <div>
          <button type="submit">Submit</button>
        </div>
      </form>
      <h1> Not a user? Signup Now! </h1>
      <button onClick={() => navigate("/register")}>Register</button>
    </div>
  );
}

const MainPage: React.FC = () => {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);

  return (< div > {isLoggedIn != false ? <Layout /> : <Login stateHandler={setIsLoggedIn} />}</div>);
}

export default MainPage
