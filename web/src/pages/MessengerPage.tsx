import { useEffect, useState, type SubmitEventHandler } from 'react';

import { getMe } from '../api/auth';
import { createChat, listChats } from '../api/chats';
import { listMessages, sendMessage } from '../api/messages';
import type { Chat } from '../types/chat';
import type { Message } from '../types/message';
import { ApiError } from '../api/client';

export default function MessengerPage() {
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
      setSelectedChatId(chat.id);
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

  return (
    <div className="messenger">

      <aside className="sidebar">
        <div className="sidebar-header">Чаты</div>
        {error && <div className="error-box" style={{ margin: '12px' }}>{error}</div>}
        <form onSubmit={handleCreateChat} style={{ padding: '12px' }}>
          <div className="field">
            <label>Username собеседника</label>
            <input
              value={peerUsername}
              onChange={(e) => setPeerUsername(e.target.value)}
              required
            />
          </div>
          <button className="btn" type="submit">Создать чат</button>
        </form>

        <div className="chat-list">
          {chats.map(chat => (
            <button
              key={chat.id}
              className={`chat-item ${selectedChatId === chat.id ? 'active' : ''}`}
              onClick={() => setSelectedChatId(chat.id)}
            >
              <div className="chat-item__name">{chat.name}</div>
              <div className="chat-item__preview">Последнее сообщение...</div>
            </button>
          )) }
        </div>
      </aside>

      <section className="chat-panel">
        {!selectedChatId ? (
          <div className="empty-state">Выбери чат слева</div>
        ) : (
          <>
            <div className="chat-header">
              {selectedChat?.name ?? 'Чат'}
            </div>
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
              <button className="btn" type="submit">→</button>
            </form>
          </>
        )}
      </section>
    </div>
  );
}