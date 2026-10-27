import { useState, type SubmitEventHandler } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { login } from '../api/auth';
import { ApiError } from '../api/client'

export default function LoginPage() {

  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);


  const handleSubmit: SubmitEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();
    setError('');


    try {
      setLoading(true);

      const data = await login({email, password});
      localStorage.setItem('token', data.token);
      navigate('/');
    } catch (error) {
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError('Login Error. Try again');
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h1>Вход</h1>
        <p>Войди в аккаунт, чтобы открыть чаты</p>

        {error && <div className="error-box">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="field">
            <label htmlFor="email">Email</label>
            <input
            id="email"
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            />
          </div>

          <div className="field">
            <label htmlFor="password">Пароль</label>
            <input 
            id="password" 
            type="password" 
            required 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <button 
          className="btn btn-primary" 
          type="submit"
          disabled={loading}
          >
            {loading ? 'Вход...' : 'Войти'}
          </button>
        </form>

        <div className="auth-link">
          Нет аккаунта? <Link to="/register">Регистрация</Link>
        </div>
      </div>
    </div>
  );
}