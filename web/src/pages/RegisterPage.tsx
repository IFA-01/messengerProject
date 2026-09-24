import { useState, type SubmitEventHandler } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { register } from '../api/auth';
import { ApiError } from '../api/client'

export default function RegisterPage() {
  const navigate = useNavigate();

  const [username, setUsername] = useState('');
  const [nickname, setNickname] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);


  const handleSubmit: SubmitEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();
    setError('');


    try {
      setLoading(true);
      await register({username, nickname, email, password});
      navigate('/login');
    } catch (error) {
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError('Register Error. Try again');
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h1>Регистрация</h1>
        <p>Создай аккаунт, чтобы открыть чаты</p>

        {error && <div className="error-box">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="field">
            <label htmlFor="username">Имя пользователя</label>
            <input
            id="username"
            type="text"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            />
          </div>

          <div className="field">
            <label htmlFor="nickname">Никнейм</label>
            <input
            id="nickname"
            type="text"
            required
            value={nickname}
            onChange={(e) => setNickname(e.target.value)}
            />
          </div>

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
            {loading ? 'Регистрация...' : 'Зарегистрироваться'}
          </button>
        </form>

        <div className="auth-link">
          Уже есть аккаунт? <Link to="/login">Вход</Link>
        </div>
      </div>
    </div>
  );
}