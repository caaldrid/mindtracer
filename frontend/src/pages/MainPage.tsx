import React, { useState } from "react";
import type { Dispatch, SetStateAction } from "react"

interface loginProp {
  stateHandler: Dispatch<SetStateAction<string | null>>
}

const Login: React.FC<loginProp> = ({ stateHandler }: loginProp) => {
  const [username, setUserName] = useState<string>();
  const [password, setPassword] = useState<string>();

  const handleLogin = async () => {
    if (username && password) {
      stateHandler(`${username} - ${password}`)
    } else {
      stateHandler(null)
    }
  }

  return (
    <div className="login-form">
      <h1> Login to continue </h1>
      <form onSubmit={handleLogin}>
        <label>
          <p>Username</p>
          <input type="text" onChange={e => setUserName(e.target.value)} required />
        </label>
        <label>
          <p>Password</p>
          <input type="password" onChange={e => setPassword(e.target.value)} required />
        </label>
        <div>
          <button type="submit">Submit</button>
        </div>
      </form>
    </div>
  );
}

const MainPage: React.FC = () => {
  const [token, setToken] = useState<string | null>(null);

  return (< div > {token != null ? <div>Yes</div> : <Login stateHandler={setToken} />}</div>);
}

export default MainPage
