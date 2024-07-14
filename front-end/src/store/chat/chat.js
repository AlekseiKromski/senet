import {createSlice} from '@reduxjs/toolkit';
import {axiosInstance} from "../axios";

const chatSlice = createSlice({
        name: 'chat',
        initialState: {
            currentChat: null,
        },
        reducers: {
            setChat(state, chat) {
                state.currentChat = chat.payload
            }
        }
    }
);

// this is for dispatch
export const {setChat} = chatSlice.actions;

// this is for configureStore
export default chatSlice.reducer;