import React, { Component, PropTypes } from 'react';
import styles from './Goo.scss';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import Paper from 'material-ui/Paper';
import RaisedButton from 'material-ui/RaisedButton';
import Divider from 'material-ui/Divider';
import * as GooActions from '../actions';
import TextField from 'material-ui/TextField';

@connect(state => ({
  goo: state.goo,
}))
export default class Goo extends Component {
  renderShortened (){
    if (!this.props.goo.rows.length)
      return;
    return (
      <Paper className={styles.urlsPaper} zDepth={1} >
        <div className={styles.urlsRow}>
          <div className={styles.urlsLong}>Long URL</div>
          <div className={styles.urlsShort}>Short URL</div>
        </div>
        {this.props.goo.rows.map(function(row,i) {
          return  <div key={i} className={styles.urlsRow}>
                    <Divider />
                    <div className={styles.urlsLong} >
                      <a href={row[1]} target="_blank" >{row[1]}</a>
                    </div>
                    <div className={styles.urlsShort} >
                      <a href={row[0]} target="_blank" >{row[0]}</a>
                    </div>
                  </div>;
        })}
      </Paper>
    );
  }
  onSubmit (){
    let url = this.refs.urlInput.input.refs.input.value;
    this.refs.urlInput.input.refs.input.value = "";
    if(!url) return;
    this.sendUrlForShort(url);
  }
  onInput (e){
    switch (e.which){
      case 13:
        e.preventDefault();
        this.onSubmit();
        break;
      default:
        return;
    }
  }
  render () {
    const actions = bindActionCreators(GooActions, this.props.dispatch);
    this.sendUrlForShort = actions.sendUrlForShort;
    return (
      <MuiThemeProvider >
        <div className={styles.goo}>
          <Paper className={styles.header} zDepth={1}>
            <h3>Just another url shortener</h3>
          </Paper>
          <Paper className={styles.inputPaper} zDepth={1} >
            <h4>Paste your long URL here:</h4>
            <div className={styles.inputSpan}>
              <TextField
                ref="urlInput"
                multiLine={true}
                className={styles.input}
                hintText="Enter URL"
                fullWidth={true}
                onKeyDown={this.onInput.bind(this)}
                autoFocus="true"
              />
            </div>
            <RaisedButton 
              secondary={true} 
              className={styles.inputSubmit} 
              label="Shorten URL" 
              backgroundColor="#a4c639" 
              onTouchEnd={this.onSubmit.bind(this)} 
              onClick={this.onSubmit.bind(this)} />
          </Paper>
          { this.renderShortened() }
        </div>
      </MuiThemeProvider>
    );
  }
}
