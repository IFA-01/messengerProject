import { Route, Routes } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import LoginPage from './pages/LoginPage';
import MessengerPage from './pages/MessengerPage';
import RegisterPage from './pages/RegisterPage';
function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route element={<ProtectedRoute />}>
        <Route path="/" element={<MessengerPage />} />
      </Route>
    </Routes>
  );
}
export default App;