import React, { Component } from 'react';
import './channel.css';

class Channel extends Component {
  render() {
    return (
      <button className={"channel " + (this.props.active ? 'active':'')}>
        <div class="channel-title">
            {this.props.title}
        </div>
        <span class="channel-icon"/>
      </button>
    );
  }
}

export default Channel;
