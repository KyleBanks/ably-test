const Client = require('./client.js');

const API_KEY = process.env.ABLY_API_KEY;
if (!API_KEY) {
    console.error('you must set the ABLY_API_KEY environment variable');
    process.exit(1);
    return;
}

const client = new Client(API_KEY);

// Not a great solution, pretty unreliable. 
// Need to do some research on the proper way to handle cleanup in Node.
const exitHandler = () => client.leaveQuiz();
process.on('exit', exitHandler.bind(null));
process.on('SIGINT', exitHandler.bind(null));

client.joinGame();
