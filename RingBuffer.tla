---------------- MODULE RingBuffer ----------------
EXTENDS Naturals

CONSTANT BufferSize
VARIABLES head, tail, buffer

\* INVARIANT: Head must never overtake Tail (Overflow)
\* INVARIANT: Tail must never overtake Head (Underflow)
TypeInvariant ==
  /\ head \in Nat
  /\ tail \in Nat
  /\ head >= tail
  /\ head - tail <= BufferSize

Init ==
  /\ head = 0
  /\ tail = 0
  /\ buffer = [i \in 0..BufferSize-1 |-> 0]

Producer ==
  /\ head < 50
  /\ head - tail < BufferSize   \* Guard: Only write if space exists
  /\ head' = head + 1           \* Update head
  /\ UNCHANGED <<tail, buffer>>

Consumer ==
  /\ head > tail                \* Guard: Only read if data exists
  /\ tail' = tail + 1           \* Update tail
  /\ UNCHANGED <<head, buffer>>

Next == Producer \/ Consumer

Spec == Init /\ [][Next]_<<head, tail, buffer>>
==============================================================================