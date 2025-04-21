## Compatibility Predictor

### Problem Overview

We need your help! Datahouse is looking to add new members to its team but needs your
amazing development skills to make the right decision. Your task, should you choose to accept it,
will be to build an application that takes an input: an array of applicants and an array of team
members, and produces an output: an array of applicants with their respective compatibility
score. How the compatibility score is generated is up to you. Help us make the right decision!

### Initial Approach

Initially, I didn't know what makes an applicant suited or "fit" for a team, so I decided to do a research.
I came across <a target="_blank" href="https://www.qic-wd.org/umbrella-summary/employee-fit">Employee Fit</a> and I learned that
compatability is not solely about just being similar to the team (vibe), but also about contributing on what the current team lacks.
This concept is known in organizational psychology as Supplementary fit (having the same vibe with the team) and Complementary Fit (filling the gaps of the team) (Kristof-Brown et al., 2005).

### Implementation

With this in mind, I wanted to use both approaches to implement my `Compatibility Predictor`.

I came accross a <a target="_blank" href="https://stackoverflow.com/questions/10561700/skill-matching-algorithm#:~:text=A%20more%20accurate%20method%20for,coordinate%20to%20the%20passed%20in">Stack Overflow forum</a> which discusses skill matching algorithm that uses Euclidean Distance. And so, I incorporated this to the first part of my implementation (Supplementary fit) by calculating the Average Fit Score of each applicant using Euclidean Distance.

To calculate `Average Fit Score`, I first calculated the team's centroid (the average values for each attribute across all team members) to where the Euclidean Distance (how far an applicant's attributes are from the team's centroid) is going to be calculated from. After getting the Euclidean Distances of each applicant, I normalized the results in order to get an easier set of data to work with by using `1 - distance / max distance`.

Just like what was mentioned above, finding compatibility does not entirely rely on Supplementary Fit but also Complementary Fit. Therefore, I also did gap analysis in order to find out how much can the applicant contribute to the current team (or their `Gap-Filler Scores`) as part of determining their compatibility with the team.

To calculate their `Gap-Filler Score`, I first calculated the gap from the current team's centroid, 1 being perfect or no gap to be filled. From this, I calculated how much each applicant can contribute to the team based on their attributes and the current team's gap per attribute.

Finally, I weighted both factors (Average Fit Score and Gap-Filler Score) to determine the final score of each applicant in terms of their compatibility with the team. I weighted Average Fit Score `80%` and Gap-Filler Score `20%` because I believe that being able to "vibe" with the team allows better workflow and work dynamics, however, it's also important to bring something new to the table, especially if working in a tech team environment, where things are constantly innovating.

### Programming Language

I decided to use GO (or GoLang) to implement this project because I became more familiar of using GO to handle JSON (the input and output type for this project) from my previous summer internship, where I had to process JSON inputs and return outputs as JSON while working on backend projects. I believe using GO would give me the less traction in accomplishing this task.

### How to run the program

In order to run this program, you first need to install Go in your system. You can follow the instructions on how to install Go in this <a target="_blank" href="https://go.dev/doc/install">documentation</a> they provided.

Once you installed GO, this project directory already includes `input.json` file which contains the sample JSON from the project description.

To run the program, make sure you're inside the project directory, then open a terminal, and run `go run main.go`.
