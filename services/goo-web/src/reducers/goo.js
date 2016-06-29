import * as types from '../constants/ActionTypes';
import omit from 'lodash/object/omit';
import assign from 'lodash/object/assign';
import mapValues from 'lodash/object/mapValues';


const initialState = {
  rows : []
};

export default function goo(state = initialState, action) {
  switch (action.type) {
    case types.SEND_URL_FOR_SHORT:
      return state;
    case types.RECEIVE_SHORTEN_URL:
      return {
        ...state,
        rows : [[action.short,action.long],...state.rows]
      }
    default:
      return state;
  }
}
