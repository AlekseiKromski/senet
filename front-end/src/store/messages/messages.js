import {createSlice} from "@reduxjs/toolkit";

const messagesSlice = createSlice({
        name: 'messages',
        initialState: {
            messages: {},
        },
        reducers: {
            setMessages: (state, p) => {
                let m = p.payload
                if (m.reverse) {
                    state.messages[m.cid] = m.messages.reverse()
                } else {
                    state.messages[m.cid] = m.messages
                }
            },
            setMessage: (state, p) => {
                let m = p.payload
                if (!state.messages[m.cid]) {
                    state.messages[m.cid] = []
                }
                state.messages[m.cid].push(m.message)
            }
        },
    }
);
export const {setMessages, setMessage} = messagesSlice.actions;

// this is for configureStore
export default messagesSlice.reducer;