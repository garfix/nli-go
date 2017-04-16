<?php
/**
 * Calls nli-go-suggest and echoes the result in JSON.
 */

$query = isset($_REQUEST['query']) ? $_REQUEST['query'] : "";
$tokens = explode(',', $query);
$sentence = implode(' ', $tokens);

$command = __DIR__ . '/../cli/nli-go-suggest "' . $sentence . '"';

// execute the Go command "nli-go-suggest"
exec($command, $output);

$suggests = [];
$error = false;
$errorLines = [];

if (count($output) == 1) {
    $suggests = json_decode($output[0], false);
} else {
    $error = true;
    $errorLines = $output;
}

header('content-type: application/json');
echo json_encode([
    'error' => $error,
    'errorLines' => $errorLines,
    'suggests' => $suggests,
], JSON_PRETTY_PRINT);
