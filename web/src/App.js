import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import InputBase from '@material-ui/core/InputBase';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import SearchIcon from '@material-ui/icons/Search';
import CloseIcon from '@material-ui/icons/Close'
import axios from 'axios';

import CodeMirror from 'codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/elegant.css';
import 'codemirror/mode/javascript/javascript.js';

import './App.css';

const styles = {
  root: {
    padding: '2px 4px',
    display: 'flex',
    alignItems: 'center',
    width: 600,
  },
  input: {
    marginLeft: 8,
    flex: 1,
  },
  iconButton: {
    padding: 10,
  },
  divider: {
    width: 1,
    height: 28,
    margin: 4,
  },
};

class CustomizedInputBase extends Component {
  state = {
    ip: "",
  }
  onSearch() {
    const { onSearch } = this.props;
    onSearch(this.state.ip);
  }
  render() {
    const { classes } = this.props;
    return (
      <Paper className={classes.root} elevation={1}>
        <IconButton className={classes.iconButton} aria-label="Menu">
          <MenuIcon />
        </IconButton>
        <InputBase
          className={classes.input}
          placeholder="Input you ip address"
          value={this.state.ip}
          onChange={(e) => this.setState({ip: e.target.value})}
        />
        <IconButton
          className={classes.iconButton}
          aria-label="Clear"
          onClick={() => this.setState({
            ip: '',
          })}
        >
          <CloseIcon />
        </IconButton>
        <Divider className={classes.divider} />
        <IconButton
          className={classes.iconButton}
          aria-label="Search"
          color="primary"
          onClick={() => this.onSearch()}
        >
          <SearchIcon />
        </IconButton>
      </Paper>
    );
  }
}

CustomizedInputBase.propTypes = {
  classes: PropTypes.object.isRequired,
  onSearch: PropTypes.func.isRequired,
};

// IP输入框
const IPLocationSearch = withStyles(styles)(CustomizedInputBase);

class App extends Component {
  locationResponder = React.createRef()
  locationResponderEditor = null
  componentDidMount() {
    this.locationResponderEditor = CodeMirror.fromTextArea(this.locationResponder.current, {
      lineNumbers: true,
      theme: 'elegant',
      mode: 'javascript',
      readOnly: true, 
    });
    // 获取当前客户端的定位
    this.doSearch('127.0.0.1');
  }
  async doSearch(ip) {
    if (!ip) {
      return
    }
    const {
      doc,
    } = this.locationResponderEditor;
    doc.setValue("// Get location by ip, please wait...");
    try {
      const {
        data,
      } = await axios.get(`/ip-location/json/${ip}`);
      // const data = {
      //   "ip": "1.0.132.192",
      //   "country": "泰国",
      //   "province": "Nakhon-Ratchasima",
      //   "city": "",
      //   "isp": "TOT"
      // };
      const curl = `// curl "https://ip.aslant.site/ip-location/json/${data.ip}\n`;
      const content = JSON.stringify(data, null, 2);
      doc.setValue(curl + content);
    } catch (err) {
      doc.setValue(`// Get location by ip fail, ${err.message}`);
    }
  }
  render() {
    const iframe = (
      //           eslint-disable-next-line
      <iframe src="https://ghbtns.com/github-btn.html?user=vicanso&repo=location&type=star&count=true&size=large" frameborder="0" scrolling="0"></iframe>
    )
    return (
      <div className="location">
        <div
          className="location-star"
        >
          {iframe}
        </div>
        <div
          className="location-ip-search"
        >
          <IPLocationSearch
            onSearch={(ip) => this.doSearch(ip)}
          />
        </div>
        <div
          className="location-responder">
          <h4>IP定位信息</h4>
          <textarea
            ref={this.locationResponder}
          ></textarea>
        </div>
        <div
          className="location-copy-right"
        >&copy; 2019 Tree Xie</div>
      </div>
    );
  }
}

export default App;
