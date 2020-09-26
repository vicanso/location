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


let searchWidth = 600;
const searchPadding = 20;
// 如果获取到浏览器大小
if (window.innerWidth && (window.innerWidth - searchPadding) < searchWidth) {
  searchWidth = window.innerWidth - searchPadding;
}

const styles = {
  root: {
    padding: '2px 4px',
    display: 'flex',
    alignItems: 'center',
    width: searchWidth,
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
          onKeyUp={(e) => {
            if (e.keyCode === 0x0d) {
              this.onSearch()
            }
          }}
          onChange={(e) => this.setState({ ip: e.target.value })}
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
  state = {
    ipLocationCount: '',
  }
  componentDidMount() {
    this.locationResponderEditor = CodeMirror.fromTextArea(this.locationResponder.current, {
      lineNumbers: true,
      theme: 'elegant',
      mode: 'javascript',
      readOnly: true,
      lineWrapping: true,
    });
    // 获取当前客户端的定位
    this.doSearch('127.0.0.1');
    this.getIPLocationCount();
  }
  async getIPLocationCount() {
    try {
      const res = await axios.get("/ip-locations/count");
      this.setState({
        ipLocationCount: res.data.count.toLocaleString(),
      });
    } catch (err) {
      console.error(err);
    }
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
      } = await axios.get(`/ip-locations/json/${ip}`);
      // const data = {
      //   "ip": "1.0.132.192",
      //   "country": "泰国",
      //   "province": "Nakhon-Ratchasima",
      //   "city": "",
      //   "isp": "TOT"
      // };
      const curl = `// curl "${location.origin}/ip-locations/json/${data.ip}"\n\n`;
      const content = JSON.stringify(data, null, 2);
      doc.setValue(curl + content);
    } catch (err) {
      let {
        message,
        response,
      } = err;
      if (response && response.data && response.data.message) {
        message = response.data.message
      }
      doc.setValue(`// Get location by ip fail, ${message}`);
    }
  }
  render() {
    const {
      ipLocationCount,
    } = this.state;
    const iframe = (
      //           eslint-disable-next-line
      <iframe src="https://ghbtns.com/github-btn.html?user=vicanso&repo=location&type=star&count=true&size=large" frameBorder="0" scrolling="0"></iframe>
    )
    const startYear = 2019;
    const currentDate = new Date();
    let copyRightDate = `${startYear}`;
    if (currentDate.getFullYear() !== startYear) {
      copyRightDate += ` - ${currentDate.getFullYear()}`;
    }
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
        <p
          className="location-curl"
        >curl "https://ip.aslant.site/ip-locations/json/8.8.8.8"</p>
        <div
          className="location-responder">
          <h4>IP Location</h4>
          <textarea
            ref={this.locationResponder}
          ></textarea>
        </div>
        <p className="location-records">There are <span>{ipLocationCount || "??"}</span> ip location records.</p>
        <div
          className="location-copy-right"
        >&copy; {copyRightDate} Tree Xie</div>
      </div>
    );
  }
}

export default App;
