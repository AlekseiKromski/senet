import {createSlice} from '@reduxjs/toolkit';

const chatsSlice = createSlice({
        name: 'chats',
        initialState: {
            chats: [],
        },
        reducers: {
            setChats(state, chats) {
                state.chats = chats.payload
            },
            setLastMessage(state, lastMessage) {
                let p = lastMessage.payload
                state.chats.find(c => c.id === p.chatId).last_message = p.message
            }
        }
    }
);

// this is for dispatch
export const {setChats, setLastMessage} = chatsSlice.actions;

// this is for configureStore
export default chatsSlice.reducer;