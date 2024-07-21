import DefaultStyle from "./default.module.css"
export default function Default(){
    return (
        <div className={DefaultStyle.DefaultWindow}>
            <h1>🛰️ Welcome to senet</h1>
            <p>On this portal you can correspond with various users. Try to find them right now! <b>Click "Find user"</b></p>
            <h2>Info</h2>
            <ul>
                <li>🔐Available basic security options - messages encrypted by server key and store in db as PGP encrypted text</li>
                <li>🤖Try to find <b>ollama3_bot</b> - this future AI bot, for any question</li>
                <li>🎙️For any questions, try to write to <b>"admin"</b> user</li>
            </ul>
        </div>
    )
}