import React, { Component } from 'react';
import Channels from './channel/channels';
import Messages from './message/messages';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <div class="aside">
          <Channels />
        </div>
        <div class="main-content">
          <Messages />
        </div>
      </div>
    );
  }
}

export default App;
