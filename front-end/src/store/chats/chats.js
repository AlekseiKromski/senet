import {createSlice} from '@reduxjs/toolkit';
import {axiosInstance} from "../axios";

const chatsSlice = createSlice({
        name: 'chats',
        initialState: {
            chats: [],
        },
        reducers: {
            setChats(state, chats) {
                state.chats = chats.payload
            }
        }
    }
);

// this is for dispatch
export const {setChats} = chatsSlice.actions;

// this is for configureStore
export default chatsSlice.reducer;