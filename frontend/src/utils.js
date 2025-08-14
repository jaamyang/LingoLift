function formatTime(seconds) {
  const date = new Date(0);
  date.setSeconds(seconds);
  const hours = date.getUTCHours().toString().padStart(2, "0");
  const minutes = date.getUTCMinutes().toString().padStart(2, "0");
  const secs = date.getUTCSeconds().toString().padStart(2, "0");
  const ms = (seconds % 1).toFixed(3).substring(2);
  return `${hours}:${minutes}:${secs},${ms}`;
}

export function exportTranscription(segments, fileName, format = 'txt') {
  let content = '';
  const fileExtension = 'txt'; // Always use .txt extension

  if (format === 'txt') {
    content = segments.map(segment => segment.text).join('\n\n');
  } else if (format === 'srt') {
    content = segments.map((segment, index) => {
      const startTime = formatTime(segment.start);
      const endTime = formatTime(segment.end);
      return `${index + 1}\n${startTime} --> ${endTime}\n${segment.text.trim()}`;
    }).join('\n\n');
  }

  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  
  const a = document.createElement('a');
  a.href = url;
  a.download = `${fileName.replace(/\.[^/.]+$/, "")}.${fileExtension}`;
  document.body.appendChild(a);
  a.click();
  
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
} 