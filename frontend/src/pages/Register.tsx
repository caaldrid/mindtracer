import React, { useState } from "react";
import { useNavigate } from "react-router-dom"
import { backendClient } from "../data/backend";

const Register: React.FC = () => {
  const [formData, setFormData] = useState({
    username: "",
    password: "",
    email: ""
  })
  const navigate = useNavigate();

  const handleRegistration: React.FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault()

    try {
      await backendClient.post("/auth/register", formData);
      navigate("/")
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
    <div className="registration-form">
      <h1> Create your account </h1>
      <form onSubmit={handleRegistration}>
        <label>
          <p>Username</p>
          <input type="text" name="username" value={formData.username} onChange={handleChange} required />
        </label>
        <label>
          <p>Password</p>
          <input type="password" name="password" value={formData.password} onChange={handleChange} required />
        </label>
        <label>
          <p>Email</p>
          <input type="email" name="email" value={formData.email} onChange={handleChange} required />
        </label>
        <div>
          <button type="submit"> Submit </button>
        </div>
      </form>
    </div>
  )
}

export default Register
