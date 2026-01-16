/**
 * TypeGo Hybrid Demo
 * 
 * This example demonstrates using NPM packages and Go modules together.
 * - lodash: Data transformation
 * - chalk: Colored terminal output
 * - go:fmt: Go's print functions
 * - go:github.com/fatih/color: Go's color library
 */

import _ from "lodash";
import chalk from "chalk";
import { Println } from "go:fmt";
import { Red, Green, Blue, Yellow, Cyan } from "go:github.com/fatih/color";

// Sample data
const users = [
    { name: "Alice", role: "admin", score: 95 },
    { name: "Bob", role: "user", score: 82 },
    { name: "Charlie", role: "user", score: 78 },
    { name: "Diana", role: "admin", score: 91 },
    { name: "Eve", role: "user", score: 88 },
];

// Header
Println("=".repeat(50));
Cyan("TypeGo Hybrid Demo: NPM + Go Modules");
Println("=".repeat(50));
Println("");

// Use lodash for data processing
const admins = _.filter(users, { role: "admin" });
const avgScore = _.meanBy(users, "score");
const topScorer = _.maxBy(users, "score");
const sortedByScore = _.orderBy(users, ["score"], ["desc"]);

// Display results using both chalk (NPM) and fatih/color (Go)
Yellow("📊 Statistics:");
Println(`  Total users: ${users.length}`);
Println(`  Average score: ${avgScore.toFixed(1)}`);
Println("");

Green("👑 Top Scorer:");
Println(`  ${topScorer?.name} with ${topScorer?.score} points`);
Println("");

Red("🔐 Administrators:");
admins.forEach((admin) => {
    Println(`  - ${admin.name} (Score: ${admin.score})`);
});
Println("");

Blue("📋 Leaderboard:");
sortedByScore.forEach((user, index) => {
    const medal = index === 0 ? "🥇" : index === 1 ? "🥈" : index === 2 ? "🥉" : "  ";
    const line = `  ${medal} ${_.padEnd(user.name, 10)} ${user.score}`;

    // Use chalk for inline styling (NPM)
    if (index === 0) {
        console.log(chalk.bold.yellow(line));
    } else if (index < 3) {
        console.log(chalk.gray(line));
    } else {
        console.log(line);
    }
});

Println("");
Println("=".repeat(50));
console.log(chalk.green.bold("✅ Demo complete!"));
