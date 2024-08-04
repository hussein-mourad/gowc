#include <algorithm>
#include <cctype>
#include <fstream>
#include <getopt.h>
#include <iomanip>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

// Args holds the command-line flags.
struct Args {
  bool printBytes = false;
  bool printChars = false;
  bool printWords = false;
  bool printLines = false;
};

// OutputData stores statistics for a single file.
struct OutputData {
  std::string file;   // file name
  int lines = 0;      // number of lines
  int words = 0;      // number of words
  int characters = 0; // number of characters
  int bytes = 0;      // number of bytes
};

std::vector<OutputData> outputData; // Stores statistics for all files
OutputData total;                   // Accumulates totals across all files
Args args;                          // Global variable for command-line flags

int maxLinesWidth = 0, maxWordsWidth = 0, maxCharsWidth = 0,
    maxBytesWidth = 0; // Widths for column formatting

// Parse command-line flags
Args parseFlags(int argc, char *argv[]) {
  Args flags;
  int option;

  while ((option = getopt(argc, argv, "cwml")) != -1) {
    switch (option) {
    case 'c':
      flags.printBytes = true;
      break;
    case 'w':
      flags.printWords = true;
      break;
    case 'm':
      flags.printChars = true;
      break;
    case 'l':
      flags.printLines = true;
      break;
    default:
      std::cerr << "Usage: " << argv[0] << " [-c] [-w] [-m] [-l] [file...]"
                << std::endl;
      exit(EXIT_FAILURE);
    }
  }

  return flags;
}

// Open a file or return stdin if the file path is empty
std::ifstream openFile(const std::string &filePath) {
  if (filePath.empty()) {
    return std::ifstream(); // Return an empty ifstream object for stdin
  }
  return std::ifstream(filePath);
}

// Calculate statistics for a given stream and file path
void calculateStats(std::istream &input, const std::string &filePath) {
  OutputData data;
  std::string line;
  std::string word;
  std::istringstream lineStream;

  while (std::getline(input, line)) {
    ++data.lines;
    data.bytes += line.size() + 1; // Include newline character
    lineStream.clear();
    lineStream.str(line);

    while (lineStream >> word) {
      ++data.words;
    }

    data.characters += line.length();
  }

  // Handle the last line if it does not end with a newline
  if (input.eof() && !line.empty()) {
    data.characters += line.length();
    std::istringstream lineStreamLast(line);
    while (lineStreamLast >> word) {
      ++data.words;
    }
  }

  // Update totals
  total.lines += data.lines;
  total.words += data.words;
  total.characters += data.characters;
  total.bytes += data.bytes;

  // Update maximum widths for formatting
  maxLinesWidth = std::max(maxLinesWidth,
                           static_cast<int>(std::to_string(data.lines).size()));
  maxWordsWidth = std::max(maxWordsWidth,
                           static_cast<int>(std::to_string(data.words).size()));
  maxCharsWidth = std::max(
      maxCharsWidth, static_cast<int>(std::to_string(data.characters).size()));
  maxBytesWidth = std::max(maxBytesWidth,
                           static_cast<int>(std::to_string(data.bytes).size()));

  data.file = filePath;
  outputData.push_back(data);
}

// Get maximum width needed for each column based on the flags
int getMaxWidth() {
  int maxWidth = 0;
  if (args.printLines)
    maxWidth = std::max(maxWidth, maxLinesWidth);
  if (args.printWords)
    maxWidth = std::max(maxWidth, maxWordsWidth);
  if (args.printChars)
    maxWidth = std::max(maxWidth, maxCharsWidth);
  if (args.printBytes)
    maxWidth = std::max(maxWidth, maxBytesWidth);
  if (!args.printLines && !args.printWords && !args.printChars &&
      !args.printBytes) {
    // Default to lines, words, and bytes
    maxWidth = std::max({maxLinesWidth, maxWordsWidth, maxBytesWidth});
  }
  return maxWidth;
}

// Print the collected output data
void printOutput() {
  int width = getMaxWidth(); // Get the maximum width for formatting

  for (const auto &data : outputData) {
    std::ostringstream output;

    if (args.printLines) {
      output << std::setw(width) << data.lines << " ";
    }
    if (args.printWords) {
      output << std::setw(width) << data.words << " ";
    }
    if (args.printChars) {
      output << std::setw(width) << data.characters << " ";
    }
    if (args.printBytes) {
      output << std::setw(width) << data.bytes << " ";
    }
    if (!args.printLines && !args.printWords && !args.printChars &&
        !args.printBytes) {
      // Default output format if no flags are provided
      output << std::setw(width) << data.lines << " ";
      output << std::setw(width) << data.words << " ";
      output << std::setw(width) << data.bytes << " ";
    }
    if (!data.file.empty()) {
      output << data.file;
    }
    std::cout << output.str() << std::endl;
  }
}

int main(int argc, char *argv[]) {
  args = parseFlags(argc, argv); // Parse command-line flags

  std::vector<std::string> filesPath;
  for (int i = optind; i < argc; ++i) {
    filesPath.push_back(argv[i]);
  }

  if (filesPath.empty()) {
    // No file paths provided, read from standard input
    std::ifstream stdinFile;
    calculateStats(std::cin, "");
    printOutput();
  } else {
    // Process each file path provided
    for (const auto &filePath : filesPath) {
      std::ifstream file = openFile(filePath);
      if (!file) {
        std::cerr << "Error opening file: " << filePath << std::endl;
        continue;
      }
      calculateStats(file, filePath);
    }

    if (filesPath.size() > 1) {
      // Add a total line if more than one file is processed
      OutputData totalArgs{"total", total.lines, total.words, total.characters,
                           total.bytes};
      outputData.push_back(totalArgs);
    }
    printOutput();
  }

  return 0;
}
