import {createAsyncThunk, createSlice} from "@reduxjs/toolkit";
import {axiosInstance} from "../axios";

const messagesSlice = createSlice({
        name: 'messages',
        initialState: {
            messages: {},
        },
        reducers: {
            setMessages: (state, p) => {
                let m = p.payload
                state.messages[m.cid] = m.messages.reverse()
            },
            setMessage: (state, p) => {
                let m = p.payload
                state.messages[m.cid].push(m.message)
            }
        },
    }
);
export const {setMessages, setMessage} = messagesSlice.actions;

// this is for configureStore
export default messagesSlice.reducer;