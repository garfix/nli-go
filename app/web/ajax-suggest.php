<?php
/**
 * Calls nli-go-suggest and echoes the result in JSON.
 */

$query = isset($_REQUEST['query']) ? $_REQUEST['query'] : "";
$tokens = explode(',', $query);
$sentence = implode(' ', $tokens);

$command = __DIR__ . '/../cli/nligo-suggest';
$configPath = __DIR__ . '/../../resources/relationships/config.json';
$fullCommand = sprintf('%s %s "%s"', $command, $configPath, $sentence);

// execute the Go command "nli-go-suggest"

$descriptorSpec = array(
   0 => array("pipe", "r"),  // stdin
   1 => array("pipe", "w"),  // stdout
   2 => array("pipe", "w")   // stderr
);

$process = proc_open($fullCommand, $descriptorSpec, $pipes);

$output = $error = "";

if (is_resource($process)) {
    $output = trim(stream_get_contents($pipes[1]));
    $error = trim(stream_get_contents($pipes[2]));
}

$suggests = empty($output) ? [] : explode("\n", $output);
$errorLines = empty($error) ? [] : explode("\n", $error);

header('content-type: application/json');
echo json_encode([
    'error' => count($errorLines) > 0,
    'errorLines' => $errorLines,
    'suggests' => $suggests,
], JSON_PRETTY_PRINT);
