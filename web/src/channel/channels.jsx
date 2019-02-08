import React, { Component } from 'react';
import Channel from './channel';
import './channel.css';

class Channels extends Component {
  render() {
    return (
      <div class="channels-container">
        <div class="channels">
          <div>
            <input class="channels-search" type="text" placeholder="Search channels" />
          </div>

          <Channel active={true} title={"Channel 1"} />
          <Channel active={false} title={"Channel 2"} />
        </div>
        <div class="add-channel-container">
          <button class="add-channel">
            <span class="add-channel-icon"></span>
            <span>Add channel</span>
          </button>
        </div>
      </div>
    );
  }
}

export default Channels;
