Features
=========

 - Restrictions
   - Paths must be a '/'-separated list of strings, like a URL or Unix filesystem.
   - All paths must begin with a '/'.
   - Path elements may not contain a '/'.
   - Trailing slashes are inconsequential.

 - Algorithm
    - Paths are mapped to the tree in the following way:
        - Each '/' is a Node in the tree. The root node is the leading '/'.
        - Each Node has edges to other nodes. The edges are named according to the possible path elements at that depth in the path.
        - Any Node may have an associated Leaf.  Leafs are terminals containing the data associated with the path as traversed from the root to that Node.

    - Edges are implemented as a map from the path element name to the next node in the path.
    
    - Extra_data is an optional information for every edge or leaf

