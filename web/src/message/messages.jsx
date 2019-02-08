import React, { Component } from 'react';
import './message.css';

class Messages extends Component {
  render() {
    return (
      <div class="messages-container">
        <div class="messages">
        </div>
        <div class="compose-message-container">
          <textarea class="message-area" placeholder="Compose message">
          </textarea>
        </div>
      </div>
    );
  }
}

export default Messages;
