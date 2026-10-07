import { useEffect, useState, type SubmitEventHandler } from 'react';
import { useNavigate } from 'react-router-dom';

import { getMe } from '../api/auth';
import { createChat, listChats } from '../api/chats';
import { listMessages, sendMessage } from '../api/messages';
import type { Chat } from '../types/chat';
import type { Message } from '../types/message';
import { ApiError } from '../api/client';
import { useMediaQuery } from '../hooks/useMediaQuery';
import { getAvatarLetter, getPeerName } from '../utils/chat';

const MOBILE_QUERY = '(max-width: 720px)';

export default function MessengerPage() {
  const navigate = useNavigate();
  const isMobile = useMediaQuery(MOBILE_QUERY);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const [myUserId, setMyUserId] = useState<number | null>(null);
  const [myUsername, setMyUsername] = useState('');

  const [chats, setChats] = useState<Chat[]>([]);
  const [selectedChatId, setSelectedChatId] = useState<number | null>(null);

  const [peerUsername, setPeerUsername] = useState('');

  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');

  const [error, setError] = useState('');

  useEffect(() => {
    getMe()
      .then((user) => {
        setMyUserId(user.id);
        setMyUsername(user.username ?? '');
                      })
      .catch(() => {
        localStorage.removeItem('token');
        window.location.href = '/login';
      });
  }, []);
  useEffect(() => {
    listChats()
    .then((data) => setChats(data.chats ?? []))
    .catch((err) => {
      if (err instanceof ApiError) {
        setError(err.message);
      }
    })
  }, []);
  useEffect(() => {
    if (!selectedChatId) {
      setMessages([]);
      return;
    }
    listMessages(selectedChatId)
      .then((data) => setMessages(data.messages ?? []))
      .catch(console.error);
  }, [selectedChatId]);
  const selectedChat = chats.find((c) => c.id === selectedChatId);

  const openChat = (chatId: number) => {
    setSelectedChatId(chatId);
    if (isMobile) setSidebarOpen(false);
  };

  const handleCreateChat: SubmitEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();
    setError('');

    try {
      if (!myUsername) {
        setError('Не удалось определить ваш username');
        return;
      }
      const generatedName = `${myUsername}-${peerUsername}`;
      const chat = await createChat({
        username: peerUsername,
        name: generatedName,
      });
      setChats((prev) => {
        const exists = prev.some((c) => c.id === chat.id);
        return exists ? prev : [...prev, chat];
      });
      openChat(chat.id);
      setPeerUsername('');
    } catch (error) {
      if (error instanceof ApiError) {
        setError(error.message);
      } else {
        setError('Create Chat Error. Try again');
      }
    }
  };

  const handleSendMessage: SubmitEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();
    if (!selectedChatId || !newMessage.trim()) return;
    try {
      const message = await sendMessage({
        chat_id: selectedChatId,
        content: newMessage.trim(),
      });
      setMessages((prev) => [...prev, message]);
      setNewMessage('');
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Send Message Error. Try again');
      }
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login', { replace: true });
  };

  return (
    <div
      className={`messenger ${sidebarOpen ? '' : 'sidebar-collapsed'}`}
      data-testid="messenger"
    >
      <div className="sidebar-backdrop" onClick={() => setSidebarOpen(false)} />

      <aside className="sidebar" inert={!sidebarOpen}>
        <div className="sidebar-header">
          <button
            className="btn btn-icon"
            type="button"
            aria-label="Скрыть список чатов"
            onClick={() => setSidebarOpen(false)}
          >
            ☰
          </button>
          Чаты
          <button className="btn btn-logout" type="button" onClick={handleLogout}>
            Выйти
          </button>
        </div>
        <form className="create-chat-form" onSubmit={handleCreateChat}>
          <input
            aria-label="Username собеседника"
            placeholder="Username собеседника"
            value={peerUsername}
            onChange={(e) => setPeerUsername(e.target.value)}
            required
          />
          <button className="btn btn-accent" type="submit">
            Создать
          </button>
        </form>
        {error && <div className="error-box">{error}</div>}

        <div className="chat-list">
          {chats.map(chat => (
            <button
              key={chat.id}
              type="button"
              className={`chat-item ${selectedChatId === chat.id ? 'active' : ''}`}
              onClick={() => openChat(chat.id)}
            >
              <div className="avatar">{getAvatarLetter(chat.name, myUsername)}</div>
              <div className="chat-item__body">
                <div className="chat-item__name">{getPeerName(chat.name, myUsername)}</div>
                <div className="chat-item__preview">Последнее сообщение...</div>
              </div>
            </button>
          )) }
        </div>
      </aside>

      <section className="chat-panel">
        <div className="chat-header">
          {!sidebarOpen && (
            <button
              className="btn btn-icon"
              type="button"
              aria-label="Показать список чатов"
              onClick={() => setSidebarOpen(true)}
            >
              ☰
            </button>
          )}
          {selectedChat ? (
            <>
              <div className="avatar avatar--small">
                {getAvatarLetter(selectedChat.name, myUsername)}
              </div>
              {getPeerName(selectedChat.name, myUsername)}
            </>
          ) : (
            'Сообщения'
          )}
        </div>
        {!selectedChatId ? (
          <div className="empty-state">
            {sidebarOpen ? 'Выбери чат' : 'Открой список чатов, чтобы выбрать чат'}
          </div>
        ) : (
          <>
            <div className="messages">
              {(messages ?? []).map(msg => (
                <div
                  key={msg.id}
                  className={`message ${msg.sender_id === myUserId ? 'own' : ''}`}
                >
                  {msg.content}
                </div>
              ))}
            </div>
            <form
              className="message-input-bar"
              onSubmit={handleSendMessage}
            >
              <input
                type="text"
                placeholder="Сообщение..."
                value={newMessage}
                onChange={e => setNewMessage(e.target.value)}
              />
              <button className="btn btn-accent" type="submit" aria-label="Отправить">
                ↑
              </button>
            </form>
          </>
        )}
      </section>
    </div>
  );
}