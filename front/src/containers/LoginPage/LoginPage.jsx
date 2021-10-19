import React, { useState } from 'react'
import Container from './LoginPage.styled'

function LoginPage({ authError, initSession }) {
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const handleSubmit = e => {
    e.preventDefault()
    const config = {
      login: login,
      password: password,
    }
    initSession(config)
  }
  return (
    <Container>
      <form className="login-form" onSubmit={handleSubmit}>
        <div className="form__field">
          <label htmlFor="login__username">
            <span className="label__hidden">Username</span>
          </label>
          <input
            id="login__username"
            type="text"
            name="username"
            className="form__input"
            placeholder="Username"
            value={login}
            required
            onChange={e => setLogin(e.target.value)}
          />
        </div>
        <div className="form__field">
          <label htmlFor="login__password">
            <span className="label__hidden">Username</span>
          </label>
          <input
            id="login__password"
            type="password"
            name="password"
            className="form__input"
            placeholder="Password"
            value={password}
            required
            onChange={e => setPassword(e.target.value)}
          />
        </div>
        <div className="form__field">
          <span className="form__field-error">{authError}</span>
        </div>
        <div className="form__field">
          <button type="submit">Log in</button>
        </div>
      </form>
    </Container>
  )
}
export default LoginPage
