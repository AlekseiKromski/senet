import { configureStore } from '@reduxjs/toolkit';
import applicationReducer from "./application/application";
import messagesReducer from "./messages/messages"
import chatReducer from "./chat/chat"
import chatsReducer from "./chats/chats"
import userReducer from "./user/user"

export default configureStore({
    reducer: {
        application: applicationReducer,
        messages: messagesReducer,
        chat: chatReducer,
        chats: chatsReducer,
        user: userReducer,
    },
});