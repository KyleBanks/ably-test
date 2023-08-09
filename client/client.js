const UUID = require('uuid');
const Ably = require('ably');
const readline = require('readline');

const CHANNEL_NAME_QUIZ = 'quiz';

const MESSAGE_NAME_QUESTION = 'question'
const MESSAGE_NAME_RESULTS = 'results'

module.exports = class Client {

  constructor(apiKey) {
    // TODO: prompt the user for a username instead of a UUID
    this.id = UUID.v4();

    this.apiKey = apiKey;
    this.inputPrompt = null;
  }

  async joinGame() {
    console.log('Joining...')
    this.client = new Ably.Realtime.Promise({
        key: this.apiKey, 
        clientId: this.id
    });
    await this.client.connection.once('connected');

    console.log('Connected, awaiting question...');
    this.channel = this.client.channels.get(CHANNEL_NAME_QUIZ);

    // This assumes we're allowed to join mid-game, we're not waiting
    // for a new game to start.
  
    await this.channel.subscribe(MESSAGE_NAME_QUESTION, (message) => {
      this._checkQuestionTimeout();
      this._presentQuestion(message.data)
    });

    await this.channel.subscribe(MESSAGE_NAME_RESULTS, (message) => {
      this._checkQuestionTimeout();
      this._presentResults(message.data);
    });
  
    // Join the game
    await this.channel.presence.enter()
  }

  async _presentQuestion(question) {
    this._openInputPrompt();

    let displayLines = [
      `\n\nQ: ${question.query}`
    ];
    question.answers.forEach((a) => {
      displayLines.push(`   ${a.id}: ${a.content}`);
    });
    displayLines.push('\n');

    const prompt = displayLines.join('\n')
    const answer = await new Promise(resolve => this.inputPrompt.question(prompt, answer => {
      resolve(answer);
    }));

    await this.channel.publish(question.id, answer);
    console.log('Answer Submitted');

    this._closeInputPrompt();
  }

  _checkQuestionTimeout() {
    if (this.inputPrompt == null)
        return;

    console.log('Timed out.')
    this._closeInputPrompt();
  }

  _presentResults(resultsData) {
    console.log('\n==============');
    console.log('Leaderboard:');

    const results = this._sortResults(resultsData);
    results.forEach(this._presentScore.bind(this));
    
    console.log('==============\n');
  }

  _presentScore(score) {
      const playerLabel = score.clientId == this.id ? '* Me *' : score.clientId;
      console.log(`${playerLabel} | ${score.score}`);
  }

  _sortResults(resultsData) {
    let results = Object.keys(resultsData.scores)
      .map(clientId => {
        return {
          clientId: clientId,
          score: resultsData.scores[clientId]
        };
      });
    results.sort((a, b) => b.score - a.score);

    return results;
  }

  async leaveQuiz() {
    if (this.client == null)
      return;

    await this.channel.presence.leave();
    this.client.close();
    this.client = null;
    console.log('Connection closed.');
  }

  _openInputPrompt() {
    this.inputPrompt = readline.createInterface({
      input: process.stdin,
      output: process.stdout
    });
  }

  _closeInputPrompt() {
    this.inputPrompt.close();
    this.inputPrompt = null;
  }
};
