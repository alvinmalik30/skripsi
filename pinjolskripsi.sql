PGDMP                         |            Pinjol    15.3    15.3 -    I           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            J           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            K           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            L           1262    16909    Pinjol    DATABASE     �   CREATE DATABASE "Pinjol" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'English_United States.1252';
    DROP DATABASE "Pinjol";
                postgres    false            �            1259    16976    application_handling_cost    TABLE     �   CREATE TABLE public.application_handling_cost (
    application_handling_cost_id character varying(55) NOT NULL,
    name character varying(100) NOT NULL,
    nominal numeric NOT NULL,
    unit character varying(100) NOT NULL
);
 -   DROP TABLE public.application_handling_cost;
       public         heap    postgres    false            �            1259    16922    biodata    TABLE     �  CREATE TABLE public.biodata (
    biodata_id character varying(55) NOT NULL,
    user_credential_id character varying(55) NOT NULL,
    full_name character varying(255),
    nik character varying(20),
    phone_number character varying(20),
    occupation character varying(255),
    place_of_birth character varying(255),
    date_of_birth date,
    postal_code character varying(10),
    is_eglible boolean,
    status_update boolean,
    additional_information text DEFAULT 'biodata is not updated'::text
);
    DROP TABLE public.biodata;
       public         heap    postgres    false            �            1259    16964    deposit    TABLE     �  CREATE TABLE public.deposit (
    deposito_id character varying(55) NOT NULL,
    user_credential_id character varying(55) NOT NULL,
    deposit_amount integer,
    interest_rate numeric,
    tax_rate numeric,
    duration integer,
    created_date date,
    maturity_date date,
    status boolean,
    gross_profit integer,
    tax integer,
    net_profit integer,
    total_return integer
);
    DROP TABLE public.deposit;
       public         heap    postgres    false            �            1259    16945    deposit_interest    TABLE     �   CREATE TABLE public.deposit_interest (
    deposito_interest_id character varying(55) NOT NULL,
    created_date date,
    interest_rate numeric,
    tax_rate numeric,
    duration_mounth integer NOT NULL
);
 $   DROP TABLE public.deposit_interest;
       public         heap    postgres    false            �            1259    17002    installenment_loan    TABLE     *  CREATE TABLE public.installenment_loan (
    installement_loan_id character varying(55) NOT NULL,
    loan_id character varying(55) NOT NULL,
    is_payed boolean,
    payment_installenment_cost integer,
    payment_deadline date,
    total_amount_of_dept integer NOT NULL,
    late_payment_fee_nominal numeric,
    late_payment_fee_unit character varying(55),
    late_payment_fee_day integer,
    late_payment_fee_total integer,
    payment_date date,
    status text,
    transfer_confirmation_recipt boolean,
    recipt_file character varying(55)
);
 &   DROP TABLE public.installenment_loan;
       public         heap    postgres    false            �            1259    17014    late_payment_fee    TABLE     �   CREATE TABLE public.late_payment_fee (
    late_payment_fee_id character varying(55) NOT NULL,
    name character varying(100) NOT NULL,
    nominal numeric NOT NULL,
    unit character varying(100) NOT NULL
);
 $   DROP TABLE public.late_payment_fee;
       public         heap    postgres    false            �            1259    16990    loan    TABLE     �  CREATE TABLE public.loan (
    loan_id character varying(55) NOT NULL,
    user_credential_id character varying(55) NOT NULL,
    loan_amount integer,
    loan_duration integer,
    loan_interest_rate numeric NOT NULL,
    loan_interest_nominal integer,
    total_amount_of_dept integer NOT NULL,
    application_handling_cost_nominal integer NOT NULL,
    application_handling_cost_unit character varying(55),
    loan_date_created date,
    is_payed boolean,
    status text
);
    DROP TABLE public.loan;
       public         heap    postgres    false            �            1259    16983    loan_interest    TABLE     �   CREATE TABLE public.loan_interest (
    loan_interest_id character varying(100) NOT NULL,
    duration_months integer NOT NULL,
    loan_interest_rate numeric NOT NULL
);
 !   DROP TABLE public.loan_interest;
       public         heap    postgres    false            �            1259    16935    saldo    TABLE     �   CREATE TABLE public.saldo (
    saldo_id character varying(55) NOT NULL,
    user_credential_id character varying(55) NOT NULL,
    total_saving integer
);
    DROP TABLE public.saldo;
       public         heap    postgres    false            �            1259    16952    top_up    TABLE     �  CREATE TABLE public.top_up (
    top_up_id character varying(55) NOT NULL,
    user_credential_id character varying(55) NOT NULL,
    top_up_amount integer,
    maturity_time timestamp without time zone,
    accepted_time timestamp without time zone,
    accepted_status boolean,
    status_information text,
    transfer_confirmation_recipt boolean,
    recipt_file character varying(55)
);
    DROP TABLE public.top_up;
       public         heap    postgres    false            �            1259    16910    user_credential    TABLE     x  CREATE TABLE public.user_credential (
    user_credential_id character varying(225) NOT NULL,
    username character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password character varying(225) NOT NULL,
    role character varying(50) NOT NULL,
    virtual_account_number character varying(225) DEFAULT ''::character varying,
    is_active boolean
);
 #   DROP TABLE public.user_credential;
       public         heap    postgres    false            B          0    16976    application_handling_cost 
   TABLE DATA           f   COPY public.application_handling_cost (application_handling_cost_id, name, nominal, unit) FROM stdin;
    public          postgres    false    220   �A       =          0    16922    biodata 
   TABLE DATA           �   COPY public.biodata (biodata_id, user_credential_id, full_name, nik, phone_number, occupation, place_of_birth, date_of_birth, postal_code, is_eglible, status_update, additional_information) FROM stdin;
    public          postgres    false    215   >B       A          0    16964    deposit 
   TABLE DATA           �   COPY public.deposit (deposito_id, user_credential_id, deposit_amount, interest_rate, tax_rate, duration, created_date, maturity_date, status, gross_profit, tax, net_profit, total_return) FROM stdin;
    public          postgres    false    219   �F       ?          0    16945    deposit_interest 
   TABLE DATA           x   COPY public.deposit_interest (deposito_interest_id, created_date, interest_rate, tax_rate, duration_mounth) FROM stdin;
    public          postgres    false    217   G       E          0    17002    installenment_loan 
   TABLE DATA           9  COPY public.installenment_loan (installement_loan_id, loan_id, is_payed, payment_installenment_cost, payment_deadline, total_amount_of_dept, late_payment_fee_nominal, late_payment_fee_unit, late_payment_fee_day, late_payment_fee_total, payment_date, status, transfer_confirmation_recipt, recipt_file) FROM stdin;
    public          postgres    false    223   !G       F          0    17014    late_payment_fee 
   TABLE DATA           T   COPY public.late_payment_fee (late_payment_fee_id, name, nominal, unit) FROM stdin;
    public          postgres    false    224   >G       D          0    16990    loan 
   TABLE DATA              COPY public.loan (loan_id, user_credential_id, loan_amount, loan_duration, loan_interest_rate, loan_interest_nominal, total_amount_of_dept, application_handling_cost_nominal, application_handling_cost_unit, loan_date_created, is_payed, status) FROM stdin;
    public          postgres    false    222   [G       C          0    16983    loan_interest 
   TABLE DATA           ^   COPY public.loan_interest (loan_interest_id, duration_months, loan_interest_rate) FROM stdin;
    public          postgres    false    221   xG       >          0    16935    saldo 
   TABLE DATA           K   COPY public.saldo (saldo_id, user_credential_id, total_saving) FROM stdin;
    public          postgres    false    216   �G       @          0    16952    top_up 
   TABLE DATA           �   COPY public.top_up (top_up_id, user_credential_id, top_up_amount, maturity_time, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file) FROM stdin;
    public          postgres    false    218   �H       <          0    16910    user_credential 
   TABLE DATA           �   COPY public.user_credential (user_credential_id, username, email, password, role, virtual_account_number, is_active) FROM stdin;
    public          postgres    false    214   �I       �           2606    16982 8   application_handling_cost application_handling_cost_pkey 
   CONSTRAINT     �   ALTER TABLE ONLY public.application_handling_cost
    ADD CONSTRAINT application_handling_cost_pkey PRIMARY KEY (application_handling_cost_id);
 b   ALTER TABLE ONLY public.application_handling_cost DROP CONSTRAINT application_handling_cost_pkey;
       public            postgres    false    220            �           2606    16929    biodata biodata_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.biodata
    ADD CONSTRAINT biodata_pkey PRIMARY KEY (biodata_id);
 >   ALTER TABLE ONLY public.biodata DROP CONSTRAINT biodata_pkey;
       public            postgres    false    215            �           2606    16951 &   deposit_interest deposit_interest_pkey 
   CONSTRAINT     v   ALTER TABLE ONLY public.deposit_interest
    ADD CONSTRAINT deposit_interest_pkey PRIMARY KEY (deposito_interest_id);
 P   ALTER TABLE ONLY public.deposit_interest DROP CONSTRAINT deposit_interest_pkey;
       public            postgres    false    217            �           2606    16970    deposit deposit_pkey 
   CONSTRAINT     [   ALTER TABLE ONLY public.deposit
    ADD CONSTRAINT deposit_pkey PRIMARY KEY (deposito_id);
 >   ALTER TABLE ONLY public.deposit DROP CONSTRAINT deposit_pkey;
       public            postgres    false    219            �           2606    17008 *   installenment_loan installenment_loan_pkey 
   CONSTRAINT     z   ALTER TABLE ONLY public.installenment_loan
    ADD CONSTRAINT installenment_loan_pkey PRIMARY KEY (installement_loan_id);
 T   ALTER TABLE ONLY public.installenment_loan DROP CONSTRAINT installenment_loan_pkey;
       public            postgres    false    223            �           2606    17020 &   late_payment_fee late_payment_fee_pkey 
   CONSTRAINT     u   ALTER TABLE ONLY public.late_payment_fee
    ADD CONSTRAINT late_payment_fee_pkey PRIMARY KEY (late_payment_fee_id);
 P   ALTER TABLE ONLY public.late_payment_fee DROP CONSTRAINT late_payment_fee_pkey;
       public            postgres    false    224            �           2606    16989     loan_interest loan_interest_pkey 
   CONSTRAINT     l   ALTER TABLE ONLY public.loan_interest
    ADD CONSTRAINT loan_interest_pkey PRIMARY KEY (loan_interest_id);
 J   ALTER TABLE ONLY public.loan_interest DROP CONSTRAINT loan_interest_pkey;
       public            postgres    false    221            �           2606    16996    loan loan_pkey 
   CONSTRAINT     Q   ALTER TABLE ONLY public.loan
    ADD CONSTRAINT loan_pkey PRIMARY KEY (loan_id);
 8   ALTER TABLE ONLY public.loan DROP CONSTRAINT loan_pkey;
       public            postgres    false    222            �           2606    16939    saldo saldo_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.saldo
    ADD CONSTRAINT saldo_pkey PRIMARY KEY (saldo_id);
 :   ALTER TABLE ONLY public.saldo DROP CONSTRAINT saldo_pkey;
       public            postgres    false    216            �           2606    16958    top_up top_up_pkey 
   CONSTRAINT     W   ALTER TABLE ONLY public.top_up
    ADD CONSTRAINT top_up_pkey PRIMARY KEY (top_up_id);
 <   ALTER TABLE ONLY public.top_up DROP CONSTRAINT top_up_pkey;
       public            postgres    false    218            �           2606    16921 )   user_credential user_credential_email_key 
   CONSTRAINT     e   ALTER TABLE ONLY public.user_credential
    ADD CONSTRAINT user_credential_email_key UNIQUE (email);
 S   ALTER TABLE ONLY public.user_credential DROP CONSTRAINT user_credential_email_key;
       public            postgres    false    214            �           2606    16917 $   user_credential user_credential_pkey 
   CONSTRAINT     r   ALTER TABLE ONLY public.user_credential
    ADD CONSTRAINT user_credential_pkey PRIMARY KEY (user_credential_id);
 N   ALTER TABLE ONLY public.user_credential DROP CONSTRAINT user_credential_pkey;
       public            postgres    false    214            �           2606    16919 ,   user_credential user_credential_username_key 
   CONSTRAINT     k   ALTER TABLE ONLY public.user_credential
    ADD CONSTRAINT user_credential_username_key UNIQUE (username);
 V   ALTER TABLE ONLY public.user_credential DROP CONSTRAINT user_credential_username_key;
       public            postgres    false    214            �           2606    16930 '   biodata biodata_user_credential_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.biodata
    ADD CONSTRAINT biodata_user_credential_id_fkey FOREIGN KEY (user_credential_id) REFERENCES public.user_credential(user_credential_id);
 Q   ALTER TABLE ONLY public.biodata DROP CONSTRAINT biodata_user_credential_id_fkey;
       public          postgres    false    215    214    3217            �           2606    16971 '   deposit deposit_user_credential_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.deposit
    ADD CONSTRAINT deposit_user_credential_id_fkey FOREIGN KEY (user_credential_id) REFERENCES public.user_credential(user_credential_id);
 Q   ALTER TABLE ONLY public.deposit DROP CONSTRAINT deposit_user_credential_id_fkey;
       public          postgres    false    219    214    3217            �           2606    17009 2   installenment_loan installenment_loan_loan_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.installenment_loan
    ADD CONSTRAINT installenment_loan_loan_id_fkey FOREIGN KEY (loan_id) REFERENCES public.loan(loan_id);
 \   ALTER TABLE ONLY public.installenment_loan DROP CONSTRAINT installenment_loan_loan_id_fkey;
       public          postgres    false    3235    222    223            �           2606    16997 !   loan loan_user_credential_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.loan
    ADD CONSTRAINT loan_user_credential_id_fkey FOREIGN KEY (user_credential_id) REFERENCES public.user_credential(user_credential_id);
 K   ALTER TABLE ONLY public.loan DROP CONSTRAINT loan_user_credential_id_fkey;
       public          postgres    false    214    222    3217            �           2606    16940 #   saldo saldo_user_credential_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.saldo
    ADD CONSTRAINT saldo_user_credential_id_fkey FOREIGN KEY (user_credential_id) REFERENCES public.user_credential(user_credential_id);
 M   ALTER TABLE ONLY public.saldo DROP CONSTRAINT saldo_user_credential_id_fkey;
       public          postgres    false    216    214    3217            �           2606    16959 %   top_up top_up_user_credential_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.top_up
    ADD CONSTRAINT top_up_user_credential_id_fkey FOREIGN KEY (user_credential_id) REFERENCES public.user_credential(user_credential_id);
 O   ALTER TABLE ONLY public.top_up DROP CONSTRAINT top_up_user_credential_id_fkey;
       public          postgres    false    214    3217    218            B   J   x�3L65I�HN�MJ31�5116׵L62�577L24121M1O�L�)���M���66�44 ΢҂���=... �L      =   �  x���ˮ[7EǺ_q@IQ�4�(2)t�	�J��G`�&��w')�`��>���Ž��ӹ"%��f)�QW�^M����o���4D�f��:z��٦Q-���!D%G�6���<���O��������ךw�GKX36���Ů�⚽�.)U)a�[�s�=�t�UV��;�ζ��M��>�N���8k�XW-1O�sh�� �8R�^ mX��.f��,#�����t��N���g���L�J�TՆ�kJ^��~9]>zx�Oo�o�`9'
��5���a�O�t謅��Ǖ�E�=�P$��R������騎�u�[����|�Q�w
���Q:$&�*�¯~��?���*L�[�[���,��B���U������A�;:3ڌZV�=i���fA-Ӧ�xUx��+�Ԓ�C�⳥�~���Sx��ؿVF��v�p������~>ߣA���z�NAo'��LQɜS��vD��I��ۘ���03���=��	��9���p�L,qd?ZS󎞲Ē�*k4OAw=P�Q�Q�Cdm@��h}.ưv�6
w��{[ص%Z<��
8�-K{)�u��hå'I�[�m�j63�G��&�����q��d�2%:V�y���׵r�qO_�w�`կ3u��J���� �hM5�w_(}��������#Vc��z$@e��2Qp8�hfmqHy1�lp�{t��kΊ��+]�Z[R�mX��0�@#-�6QjA�*60�c�i��$jhI��7�j�H��6�|��yta��ɰU�J(�E*�`�$x��	�H���_�k�"X��0��љ+�1{�<�$0\���Hp� �K�u�|��V�۲dV��ɭy�/��N).��5�=�B9�60�c㢊,��q�����i��8����}-��d��#�� S&F�F�|�$ƤOXz��ET3	{<=���˻����t͑Ƈ�v�_�DZ�����ЇA'\��G�#��̞�T��o�H��(mڐ�!�AZFD��t���#�F �A�wSz$^� �¤�ο���:���#H��/w	7�`��Z9��y�(ll���w�>~�#�#��=a�w��n���jF�cC�����L�O��L��;@Y3bp%��[@��)�)�����r-�o���S�µ}zaa�/�6�Th��Ύ�\^�)M�Oi>u��������      A      x������ � �      ?      x������ � �      E      x������ � �      F      x������ � �      D      x������ � �      C   P   x���� �7�ő����迄����K��;M'����z���⥉c�j�{�1�r��cl�=�񄠙����'	�      >   �   x�5�An 1�ٿ�"Y�/�/��o*��|�uYc���_�/�n���)c�B�00��U'�O�pi��n�*@1,:���h*�����n�wbXp2�cN��T�K��+[�A��V{���ƺ���_%��b����}��L�{�{�i4��j@$�z�c�w�X���-����x��DE      @   �   x����j�0D��W��J�[�����!�nB6���M�P�bo3Y�S����'Ap���e�.�tlj�^�,xP�E�U�Z*"���a��ð$`��=s�C<u��������1q2����u/�a�����۾�����Q�	֭<���8� �:�Юש��;��
5�|��D���-�	,q蠟�v�^�90�K��T���>_'{_���]�ݷ|�ղ٭輖#Լ�M�|9^�      <   �  x���GwbI�ת���D�ͮ�F8Bg6i�yx'��O-U����l �����7n$�Ά�lu��ˍ]�
K;�
���_��E�_��m�ۜ��]�z����M�]g�7�µ�}XMU���Te>�̦�V��m��w��n����� ������8��#΃C6E��JZb⸉����V����}�h��c?������L�q�o�ټŰ�>�V��<��}��x��Juui�m"�tn�w*j�g�% !d�H*l����h�4QkHJaO9�BQd�w�&"TPXK+�6�-fng��忛����X����qa ��i���8�G�����a���]��N��AI����:���X�uQ�)�V&i�B�!�p�j��c�	Y'�^i�-b�T$�Ƽ�6;�V�������9��ū}�	>�ǆ�dY|���a1b�ڵN�8ꏦ���q[$��R�%���JE1F�5�ș���RD�-3��[�M�	�<��dX��c�T[/��~�L�����V(ȸ���0�-�Z�h�:��f>�
��k�RpO��u�/����<*�?��K�B̂U���A}?`/� �@� v�X��~�g��,�����y'$��a
o�C�:�4�)�@#����ߧ%w�o�q��r�*��~[̶���Ӽ<��=��Ʋ(��.kNjE�P=xz�'�mnf�ګ �Q�'2��p�h,�W)P̬q�@��h�c�!"�����؟����[��~\<�v�I��]�aRhͧ��-����!m�v����}U|z'��/q��� ��P� ��L�'gU>*�K���R伅]�-$���kY0�hߖ�V���`�Y�����_ļ5���I��|ל�6���o�4Q�b���շ�H�1�=G"D>a8��	�,�H�U��á���QI�1�1����7�)����m-n�����G�i�Z:�˥u����ʰj����~ٽ�OU�6j�:;6������e�K�% ��#�U.a������!yR�<�[�8�HyK�L	�>���Ɣܐ���^[��5o����ݾ9,��\Tu�;��\��{g�m7��.\t瓔J�|.��H�͂�:R�U,'�>���(x����w���Q"Y�w_�%_����*�Z���Ǖ��S=t��8�J�$����;����aۡ�;n]�/x�)wF#�h{��1�8k!�iT �D$X
���7!��Uz�@B�ؐ�\|}ߠf-�Z�l�8[�^�K�\�E�S��3[5���=½��l�C�,?N6y SI����B���4!ZF��R9'7�F�<�sE��$A�F�al��� �u�f���~lnH{C�3�t�ʲ���XV�Ǝ�%�����<�CR�����ޝ_R_l_1����Z��!WaDB�A�����s ��,xT�z0�,����Kpt���~��m���67��^��N��Z������Z.�!4��o�����m���j<���/��B)$m���1yA�	��Q�}�̔u����\� � hH*��X_�_����tu;rY!\�8܋����I�Y�tà`���c�Fu~2�M�Q��M�����M���A���"�1] ��Ri��1������!.��.p")�k�'6���6 *�ª�E�<�&�Ċ�^�s|�&���������\6�nU�g��x��o7����@L��s�]�E����ֲ�S�-N	�!J�����/Lb�m ~?�����N��bV�oO�S�<�U�J'[uJ�A1�.�����>����(쪾��_א�i�Xϭ�	��SC����	�}�XIYBX�<��C	�����$�Ն�x��o���Z<����>��*��qw�u���-���ʙ,�ĲI/ꭹh��,վݱe���#��<�$�����M���H
7�V�G
g�gD���y!�i�F�J��ހ�����eQ�"��[���t�n�z����k����b����������,�+H/����0Ő���N'(>��~���ױ�     