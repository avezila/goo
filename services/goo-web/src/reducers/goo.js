import * as types from '../constants/ActionTypes';


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
