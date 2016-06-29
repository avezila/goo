import React, { Component } from 'react';

import  { createStore as initStore, combineReducers, compose, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import thunkMiddleware from 'redux-thunk';
import persistState from 'redux-localstorage'



import injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();


import Goo from "./Goo"
import * as reducers from '../reducers';


function lsSlicer (paths) {
  return (state) => {
    console.log(state)
    let subset = {
      goo : {
        rows  : state.goo.rows
      }
    }
    return subset
  }
}

const enhancer = compose(
  applyMiddleware(thunkMiddleware),
  persistState("redux",{slicer:lsSlicer}),
)

//let createStore = compose(
//)();

const reducer = combineReducers(reducers);
const store = initStore(
  reducer, enhancer
);


export default class App extends Component {
  render () {   
    return (
      <div>
        <Provider store={store} >
          <Goo />
        </Provider>
      </div>
    );
  }
}
