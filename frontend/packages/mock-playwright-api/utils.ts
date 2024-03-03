export function spellWord(word: string): string[] {
  return word.split('').map((letter, index, arr) => {
    const out = arr.slice(0, index + 1);
    const rv = out.length === 0 ? [letter] : out;
    return out.join('');
  });
}
