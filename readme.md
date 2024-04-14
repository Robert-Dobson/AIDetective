# AI Detective

For the 24 hour hackathon at the University of Bath (Bath Hack 2024), we made a social deduction game where a dective needs to deduce who is human and who is AI. 

With AI scams and deepfakes becoming more and more worrying, this game challenges your belief that you would be able to identify what is generated by an AI. AI Dectective is a creative game that helps educate players the dangers of AI. This is a web application that can be played on both your phones and computers.

[ Add GIF ] 

## How the Game Works
Throughout the game there is a single detective. The goal of the detective is to successfully identify the humans and eliminate all of them. There are also up to 4 humans, the goal of these players is to blend in with the AI. If at least one human survives the rounds and avoids being identified by the dectective, they will win. The last role is the AI which, as the name suggests, are powered by AI. Their goal is to act as closely to humans as possible. 

Each round detective can provide a prompt to all the players (humans and AI) for which they must respond. Using the responses given, the detective then attempts to choose a human to eliminate. At the end of the round, the true identity of the eliminated player will be revealed, and depending on if any Humans remain, they will begin a new round. 

[ Add photos]

## How we Developed it
- Front-end was written in HMTL, CSS and JavaScript
- Back-end was written in Golang, taking advantage of the [Melody library](https://github.com/olahol/melody/tree/master)
- The AI responses were generated by using an OpenAI GPT3.5 Turbo API Key 

